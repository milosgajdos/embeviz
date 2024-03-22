import ReactECharts from "echarts-for-react";
import "echarts-gl";
import { getChartOption } from "./options";

export default function ECharts({ isLoading, embeddings }) {
  return (
    <div id="echarts">
      <EChart dim="3D" isLoading={isLoading} embeddings={embeddings["3D"]} />
      <EChart dim="2D" isLoading={isLoading} embeddings={embeddings["2D"]} />
    </div>
  );
}

function EChart({ dim, isLoading, embeddings }) {
  // TODO: should this be a hook?
  let option = getChartOption(dim, embeddings);

  return (
    <ReactECharts
      notMerge
      showLoading={isLoading}
      option={option}
      className="echart-container"
    />
  );
}
