import ReactECharts from "echarts-for-react";
import "echarts-gl";
import { getChartOption } from "./options";

export default function EChart({ dim, isLoading, series, styling }) {
  // TODO: should this be a hook?
  const option = getChartOption(dim, series);
  return (
    <ReactECharts
      notMerge
      showLoading={isLoading}
      option={option}
      style={styling} // TODO: fix this styling mess
    />
  );
}
