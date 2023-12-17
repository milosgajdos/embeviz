import ReactECharts from "echarts-for-react";
import "echarts-gl";
import { getChartOption } from "./options";

export default function EChart({ name, dim, isLoading, embeddings, styling }) {
  // TODO: should this be a hook?
  const option = getChartOption(name, dim, embeddings);
  return (
    <ReactECharts
      notMerge
      showLoading={isLoading}
      option={option}
      style={styling} // TODO: fix this styling mess
    />
  );
}
