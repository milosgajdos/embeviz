import { Form, useLoaderData, useNavigation, redirect } from "react-router-dom";
import { useState } from "react";
import {
  getProvider,
  getProviderProjections,
  updateData,
  deleteData,
} from "../lib/embeddings";
import EChart from "../components/echart/echart";

async function getReqData(request) {
  const formData = await request.formData();
  return Object.fromEntries(formData);
}

export async function action({ request, params }) {
  switch (request.method) {
    case "POST": {
      await updateData(params.uid, await getReqData(request));
      break;
    }
    case "DELETE": {
      const isConfirmed = window.confirm(
        "Are you sure you want to delete the data?",
      );
      if (isConfirmed) {
        await deleteData(params.uid);
      }
      break;
    }
    default: {
      throw new Response("Unsupported operation", { status: 405 });
    }
  }
  return { ok: true };
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
      <UpdateForm />
    </>
  );
}

export function UpdateForm() {
  const [projection, setProjection] = useState("pca");

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
          <input id="label" name="label" placeholder="Label" />
          <textarea
            id="text"
            name="text"
            placeholder="Text"
            rows="10"
            cols="80"
            wrap="soft"
            required
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
          <button type="submit">Embed</button>
          <button type="button" id="delete-btn" onClick={handleClearFields}>
            Clear
          </button>
        </div>
      </Form>
      <div id="embed-buttons">
        <Form method="delete">
          <button type="submit" id="delete-btn">
            Drop Data
          </button>
        </Form>
      </div>
    </>
  );
}
