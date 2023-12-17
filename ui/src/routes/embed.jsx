import { Form, useLoaderData, useNavigation, redirect } from "react-router-dom";
import ReactECharts from "echarts-for-react";
import "echarts-gl";
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
            <EChart3D
              isLoading={navigation.state === "loading"}
              series={provider.data["3D"]}
            />
            <EChart2D
              isLoading={navigation.state === "loading"}
              series={provider.data["2D"]}
            />
          </div>
        </div>
      </div>
      <UpdateDataForm />
    </>
  );
}

function EChart3D({ isLoading, series }) {
  const option = {
    animation: true,
    legend: { show: true, type: "" },
    grid3D: {},
    xAxis3D: {},
    yAxis3D: {},
    zAxis3D: {},
    toolbox: {
      show: true,
      orient: "horizontal",
      left: "right",
      feature: {
        saveAsImage: {
          show: true,
          title: "Save as image",
        },
        brush: null,
        restore: {
          show: true,
          title: "Reset",
        },
      },
    },
    tooltip: {
      show: true,
      formatter: "{a}",
    },

    series: [
      {
        name: "Sample 3D data",
        type: "scatter3D",
        symbolSize: 5,
        smooth: false,
        connectNulls: false,
        showSymbol: false,
        waveAnimation: false,
        coordinateSystem: "cartesian3D",
        renderLabelForZeroData: false,
        data: series,
        itemStyle: {
          opacity: 1,
        },
      },
    ],
  };
  return (
    <ReactECharts
      notMerge
      showLoading={isLoading}
      option={option}
      style={{ height: 300, width: 300 }}
    />
  );
}

function EChart2D({ isLoading, series }) {
  const option = {
    animation: true,
    legend: { show: true, type: "" },
    xAxis: [{}],
    yAxis: [{}],
    toolbox: {
      show: true,
      orient: "horizontal",
      left: "right",
      feature: {
        saveAsImage: {
          show: true,
          title: "Save as image",
        },
        brush: null,
        restore: {
          show: true,
          title: "Reset",
        },
      },
    },
    tooltip: {
      show: true,
      formatter: "{a}",
    },

    series: [
      {
        name: "Sample 2D data",
        type: "scatter",
        symbolSize: 5,
        smooth: false,
        connectNulls: false,
        showSymbol: false,
        waveAnimation: false,
        renderLabelForZeroData: false,
        data: series,
        itemStyle: {
          opacity: 1,
        },
      },
    ],
  };
  return (
    <ReactECharts
      notMerge
      showLoading={isLoading}
      option={option}
      style={{ height: 300, width: 300 }}
    />
  );
}

export function UpdateDataForm() {
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
