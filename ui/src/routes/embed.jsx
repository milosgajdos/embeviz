import {
  Form,
  useLoaderData,
  useNavigation,
  redirect,
  useParams,
  useRevalidator,
  json,
} from "react-router-dom";
import { useState } from "react";
import {
  getProvider,
  getProviderProjections,
  embedData,
  computeData,
  deleteData,
} from "../lib/embeddings";
import EChart from "../components/echart/echart";

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

  // NOTE: we need to "revalidate" the parent component
  // if we delete the data so the charts are rerendered.
  // I tried using state but that somehow never triggers render.
  function handleDataDeleted() {
    revalidator.revalidate();
  }

  return (
    <>
      <div id="embed">
        <div>
          <h1>{provider.name ? <> {provider.name}</> : <i>No Name</i>} </h1>
          {provider.description && <p>{provider.description}</p>}
          <div id="charts">
            <EChart
              name="3D projections"
              dim="3D"
              isLoading={navigation.state === "loading"}
              embeddings={embeddings["3D"]}
              styling={{ height: 300, width: 300 }}
            />
            <EChart
              name="2D projections"
              dim="2D"
              isLoading={navigation.state === "loading"}
              embeddings={embeddings["2D"]}
              styling={{ height: 300, width: 300 }}
            />
          </div>
        </div>
      </div>
      <UpdateForm onDataDeleted={handleDataDeleted} />
    </>
  );
}

export function UpdateForm({ onDataDeleted }) {
  let params = useParams();
  const [projection, setProjection] = useState("pca");

  async function handleDeleteData() {
    const isConfirmed = window.confirm(
      "Are you sure you want to delete the data?",
    );
    if (isConfirmed) {
      try {
        await deleteData(params.uid);
        onDataDeleted();
      } catch (error) {
        console.error("Error deleting data:", error);
      }
    }
  }

  function handleProjectionChange(e) {
    setProjection(e.target.value);
  }

  function handleClearFields() {
    document.getElementById("label").value = "";
    document.getElementById("text").value = "";
  }

  return (
    <>
      <Form method="post" id="embed-form">
        <div id="embed-form-text-options">
          <input id="label" name="label" placeholder="Label (Optional)" />
          <textarea
            id="text"
            name="text"
            placeholder="Text"
            rows="10"
            cols="80"
            wrap="soft"
          ></textarea>
        </div>
        <div id="embed-options">
          <fieldset>
            <legend>Projection</legend>
            <div>
              <input
                type="radio"
                id="pca"
                name="projection"
                value="pca"
                checked={projection === "pca"}
                onChange={handleProjectionChange}
              />
              <label htmlFor="pca"> pca</label>
            </div>
            <div>
              <input
                type="radio"
                id="tsne"
                name="projection"
                value="tsne"
                checked={projection === "tsne"}
                onChange={handleProjectionChange}
              />
              <label htmlFor="tsne"> t-SNE</label>
            </div>
          </fieldset>
          <button type="submit" name="intent" value="embed">
            Embed
          </button>
          <button type="submit" id="update-btn" name="intent" value="compute">
            Compute
          </button>
          <button type="button" id="delete-btn" onClick={handleClearFields}>
            Clear
          </button>
        </div>
      </Form>
      <div id="embed-buttons">
        <Form method="delete">
          <button type="button" id="delete-btn" onClick={handleDeleteData}>
            Drop Data
          </button>
        </Form>
      </div>
    </>
  );
}
