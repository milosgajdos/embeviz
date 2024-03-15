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
  let intent = formData.get("intent");
  const updates = Object.fromEntries(formData);

  switch (intent) {
    case "embed":
      await embedData(params.uid, updates);
      break;
    case "compute":
      await computeData(params.uid, updates);
      break;
    default:
      throw json({ message: "Invalid intent" }, { status: 400 });
  }
  return redirect(`/provider/${params.uid}`);
}

export async function loader({ params }) {
  // TODO: fetch provider and embeddings in parallel
  const provider = await getProvider(params.uid);
  if (!provider) {
    throw new Response("", {
      status: 404,
      statusText: "Not Found",
    });
  }
  const { embeddings } = (await getProviderProjections(params.uid)) ?? [];
  return { provider, embeddings };
}

export default function Embed() {
  const { provider, embeddings } = useLoaderData();
  const navigation = useNavigation();
  const revalidator = useRevalidator();
  const [isFetching, setFetching] = useState(false);

  // NOTE: we need to "revalidate" the parent component
  // if we delete the data so the charts are rerendered.
  function handleDeletion() {
    revalidator.revalidate();
  }

  function handleFetching(fetching) {
    setFetching(fetching);
  }

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
      <EmbedForm onDeletion={handleDeletion} onFetching={handleFetching} />
    </>
  );
}
