import {
  useLoaderData,
  useNavigation,
  redirect,
  useRevalidator,
  json,
} from "react-router-dom";
import { useState } from "react";
import {
  getProvider,
  getProviderProjections,
  embedData,
  computeData,
  getInputChunks,
  computeChunks,
} from "../lib/api";
import ECharts from "../components/echarts/echarts";
import EmbedForm from "../components/embed-form/embed-form";

export async function action({ request, params }) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const updates = Object.fromEntries(formData);

  console.log("intent: " + intent);

  // TODO: handle errors better:
  // Status message should inform the user WHY the failure happened
  // so the user goes back and fixes the problem if possible
  switch (intent) {
    case "embed":
      try {
        await embedData(params.uid, updates);
      } catch (error) {
        console.log(error);
        throw new Response("", {
          status: error.Status,
          statusText: "Invalid input",
        });
      }
      break;
    case "compute":
      try {
        await computeData(params.uid, updates);
      } catch (error) {
        throw new Response("", {
          status: error.Status,
          statusText: "Computing projections failed!",
        });
      }
      break;
    case "chunk":
      try {
        console.log("chunking", updates);
        await computeChunks(updates);
      } catch (error) {
        throw new Response("", {
          status: error.Status,
          statusText: "Getting chunks failed!",
        });
      }
      break;
    default:
      throw json({ message: "Invalid intent" }, { status: 400 });
  }
  return redirect(`/provider/${params.uid}`);
}

export async function loader({ params }) {
  // TODO: fetch provider and embeddings in parallel
  let provider, embeddings, chunks;
  try {
    provider = await getProvider(params.uid);
    if (!provider) {
      throw new Response("", {
        status: 404,
        statusText: "Not found!",
      });
    }
  } catch (error) {
    throw new Response("", {
      status: error.Status,
      statusText: "Failed reading provider!",
    });
  }
  try {
    ({ embeddings } = await getProviderProjections(params.uid));
  } catch (error) {
    throw new Response("", {
      status: error.Status,
      statusText: "Fetching projections failed!",
    });
  }
  try {
    ({ chunks } = await getInputChunks());
  } catch (error) {
    throw new Response("", {
      status: error.Status,
      statusText: "Fetching chunks failed!",
    });
  }

  return { provider, embeddings, chunks };
}

export default function Embed() {
  const { provider, embeddings, chunks } = useLoaderData();
  const navigation = useNavigation();
  const revalidator = useRevalidator();
  const [isFetching, setFetching] = useState(false);

  return (
    <>
      <div id="embed">
        <h1>{provider.name ? <> {provider.name}</> : <i>No Name</i>} </h1>
        {provider.description && <p>{provider.description}</p>}
        <ECharts
          isLoading={navigation.state === "loading" || isFetching}
          embeddings={embeddings}
        />
      </div>
      <EmbedForm
        // NOTE: we need to "revalidate" the parent component
        // if we drop the data so the charts are rerendered.
        onDrop={() => revalidator.revalidate()}
        onFetch={(fetching) => setFetching(fetching)}
        chunks={chunks}
      />
    </>
  );
}
