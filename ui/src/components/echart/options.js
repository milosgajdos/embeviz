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
    title: {
      text: "3D Projections",
    },
  },
  "2D": {
    animation: true,
    xAxis: [{}],
    yAxis: [{}],
    legend: defaultLegend,
    toolbox: defaultToolbox,
    tooltip: defaultToolTip,
    series: [],
    title: {
      text: "2D Projections",
    },
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

export function makeSeries(dim, data) {
  // TODO: clean this up, it's very fugly
  data = data.map((obj) => {
    const { metadata } = obj;
    const name = metadata?.label;
    const color = metadata?.color;

    if (name) {
      obj = {
        ...obj,
        name: name,
      };
    }

    if (color) {
      obj = {
        ...obj,
        itemStyle: {
          color: color,
        },
      };
    }

    return obj;
  });

  return {
    ...defaultSeriesOptions[dim],
    data: data,
  };
}

export function getChartOption(dim, embeddings) {
  return {
    ...defaultChartOptions[dim],
    series: makeSeries(dim, embeddings),
  };
}
