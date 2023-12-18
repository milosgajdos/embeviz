import { Form, useLoaderData, useNavigation, redirect } from "react-router-dom";
import { useState } from "react";
import { getProvider, updateData } from "../lib/embeddings";
import EChart from "../components/echart/echart";

export async function loader({ params }) {
  const provider = await getProvider(params.id);
  if (!provider) {
    throw new Response("", {
      status: 404,
      statusText: "Not Found",
    });
  }
  console.log(provider);
  return { provider };
}

export async function action({ request, params }) {
  const formData = await request.formData();
  const updates = Object.fromEntries(formData);
  console.log(updates);
  const newData = await updateData(params.id, updates);
  console.log(`new data: ${newData}`);
  return redirect(`/provider/${params.id}`);
}

export default function Embed() {
  const { provider } = useLoaderData();
  const navigation = useNavigation();

  return (
    <>
      <div id="embed">
        <div>
          <h1>{provider.name ? <> {provider.name}</> : <i>No Name</i>} </h1>
          {provider.description && <p>{provider.description}</p>}
          <div id="charts">
            <EChart
              name="3D series"
              dim="3D"
              isLoading={navigation.state === "loading"}
              embeddings={provider.embeddings["3D"]}
              styling={{ height: 300, width: 300 }}
            />
            <EChart
              name="2D series"
              dim="2D"
              isLoading={navigation.state === "loading"}
              embeddings={provider.embeddings["2D"]}
              styling={{ height: 300, width: 300 }}
            />
          </div>
        </div>
      </div>
      <UpdateDataForm />
    </>
  );
}

export function UpdateDataForm() {
  const [projection, setProjection] = useState("pca");

  function handleProjectionChange(e) {
    setProjection(e.target.value);
  }

  return (
    <Form action="update" method="post" id="embed-form">
      <div id="embed-form-text-options">
        <input id="meta" name="meta" placeholder="Label" />
        <textarea
          id="text"
          name="text"
          placeholder="Text"
          rows="5"
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
            <label htmlFor="tsne"> t-sne</label>
          </div>
        </fieldset>
        <button type="submit">Submit</button>
      </div>
    </Form>
  );
}
