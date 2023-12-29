const defaultToolbox = {
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
};

const defaultToolTip = {
  show: true,
  formatter: "{b} [{c}]",
};

const defaultLegend = {
  show: true,
  type: "",
};

const defaultChartOptions = {
  "3D": {
    animation: true,
    grid3D: {},
    xAxis3D: {},
    yAxis3D: {},
    zAxis3D: {},
    legend: defaultLegend,
    toolbox: defaultToolbox,
    tooltip: defaultToolTip,
    series: [],
  },
  "2D": {
    animation: true,
    xAxis: [{}],
    yAxis: [{}],
    legend: defaultLegend,
    toolbox: defaultToolbox,
    tooltip: defaultToolTip,
    series: [],
  },
};

const defaultSeries = {
  symbolSize: 5,
  smooth: false,
  connectNulls: false,
  showSymbol: false,
  waveAnimation: false,
  renderLabelForZeroData: false,
  itemStyle: { opacity: 1 },
  data: [],
};

const defaultSeriesOptions = {
  "2D": {
    ...defaultSeries,
    type: "scatter",
  },
  "3D": {
    ...defaultSeries,
    type: "scatter3D",
    coordinateSystem: "cartesian3D",
  },
};

export function makeSeries(name, dim, data) {
  // extracting the value of label and  setting it as a name
  // TODO: clean this up, it's very fugly
  data = data.map((obj) => {
    const { metadata, ...rest } = obj;
    const name = metadata?.label;
    if (name) {
      return {
        ...rest,
        name,
      };
    } else {
      return obj;
    }
  });

  return {
    ...defaultSeriesOptions[dim],
    name: name,
    data: data,
  };
}

export function getChartOption(name, dim, embeddings) {
  return {
    ...defaultChartOptions[dim],
    series: makeSeries(name, dim, embeddings),
  };
}
