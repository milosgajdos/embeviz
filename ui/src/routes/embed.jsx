import { Form, useLoaderData, redirect } from "react-router-dom";
import { useState } from "react";
import { getProvider, updateData } from "../embeddings";

export async function loader({ params }) {
  const provider = await getProvider(params.id);
  if (!provider) {
    throw new Response("", {
      status: 404,
      statusText: "Not Found",
    });
  }
  return { provider };
}

export async function action({ request, params }) {
  const formData = await request.formData();
  const updates = Object.fromEntries(formData);
  const newData = await updateData(params.id, updates);
  console.log(`new data: ${newData}`);
  return redirect(`/provider/${params.id}`);
}

export default function Embed() {
  const { provider } = useLoaderData();
  /* we should allow reset stored data*/
  return (
    <>
      <div id="embed">
        <h1>{provider.name ? <> {provider.name}</> : <i>No Name</i>} </h1>
        {provider.description && <p>{provider.description}</p>}
        <div>
          <img key={provider.avatar} src={provider.avatar || null} />
        </div>
      </div>
      <UpdateData />
    </>
  );
}

export function UpdateData() {
  const [proj, setProj] = useState("pca");

  function handleProjChange(e) {
    setProj(e.target.value);
  }

  return (
    <Form action="update" method="post" id="embed-form">
      <textarea
        id="text"
        name="text"
        placeholder="Text"
        rows="5"
        cols="80"
        wrap="soft"
        required
      ></textarea>
      <div id="embed-options">
        <fieldset>
          <legend>Projection</legend>
          <div>
            <input
              type="radio"
              id="pca"
              name="proj"
              value="pca"
              checked={proj === "pca"}
              onChange={handleProjChange}
            />
            <label htmlFor="pca"> pca</label>
          </div>
          <div>
            <input
              type="radio"
              id="tsne"
              name="proj"
              value="tsne"
              checked={proj === "tsne"}
              onChange={handleProjChange}
            />
            <label htmlFor="tsne"> t-sne</label>
          </div>
        </fieldset>
        <button type="submit">Submit</button>
      </div>
    </Form>
  );
}
