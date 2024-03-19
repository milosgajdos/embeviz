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
} from "../lib/embeddings";
import EChart from "../components/echart/echart";
import EmbedForm from "../components/embed-form/embed-form";

export async function action({ request, params }) {
  const formData = await request.formData();
  const intent = formData.get("intent");
  const updates = Object.fromEntries(formData);

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
    default:
      throw json({ message: "Invalid intent" }, { status: 400 });
  }
  return redirect(`/provider/${params.uid}`);
}

export async function loader({ params }) {
  // TODO: fetch provider and embeddings in parallel
  let provider;
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
    const { embeddings } = await getProviderProjections(params.uid);
    return { provider, embeddings };
  } catch (error) {
    throw new Response("", {
      status: error.Status,
      statusText: "Fetching projections failed!",
    });
  }
}

export default function Embed() {
  const { provider, embeddings } = useLoaderData();
  const navigation = useNavigation();
  const revalidator = useRevalidator();
  const [isFetching, setFetching] = useState(false);

  return (
    <>
      <div id="embed">
        <div>
          <h1>{provider.name ? <> {provider.name}</> : <i>No Name</i>} </h1>
          {provider.description && <p>{provider.description}</p>}
          <br />
          <div id="charts">
            <EChart
              dim="3D"
              isLoading={navigation.state === "loading" || isFetching}
              embeddings={embeddings["3D"]}
            />
            <EChart
              dim="2D"
              isLoading={navigation.state === "loading" || isFetching}
              embeddings={embeddings["2D"]}
            />
          </div>
        </div>
      </div>
      <EmbedForm
        // NOTE: we need to "revalidate" the parent component
        // if we drop the data so the charts are rerendered.
        onDrop={() => revalidator.revalidate()}
        onFetching={(fetching) => setFetching(fetching)}
      />
    </>
  );
}
