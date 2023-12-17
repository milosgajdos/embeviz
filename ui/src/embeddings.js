import localforage from "localforage";
import { matchSorter } from "match-sorter";
import sortBy from "sort-by";
import { makeSeries, makeSeriesItem, addDataItem } from "./charts/options";

const globalData = new Map([
  [
    "openai",
    {
      name: "OpenAI",
      description: "OpenAI Embeddings",
      embeddings: {
        "2D": [
          [0.5, 0.6],
          [0.1, 0.2],
          [0.4, 0.1],
        ],
        "3D": [
          [0.5, 0.6, 0.3],
          [0.1, 0.2, 0.1],
          [0.4, 0.1, 0.2],
        ],
      },
    },
  ],
  [
    "vertexai",
    {
      name: "VertexAI",
      description: "VertexAI Embeddings",
      embeddings: {
        "2D": [
          [0.43, 0.77],
          [0.21, 0.33],
          [0.32, 0.68],
        ],
        "3D": [
          [0.43, 0.77, 0.453],
          [0.21, 0.33, 0.51],
          [0.42, 0.78, 0.62],
        ],
      },
    },
  ],
]);

export async function getProviders(query) {
  await fakeNetwork(`getProviders:${query}`);
  let providers = [];
  for (const [id, obj] of globalData) {
    providers.push({ id: id, name: obj.name });
  }
  if (!providers) return [];
  if (query) {
    providers = matchSorter(providers, query, { keys: ["name"] });
  }
  return providers.sort(sortBy("name"));
}

export async function getProvider(id) {
  await fakeNetwork(`provider:${id}`);
  let provider = globalData.get(id);
  return provider ?? null;
}

export async function updateData(id, updates) {
  const { meta, text, projection } = updates;
  await fakeNetwork();
  let provider = globalData.get(id);
  if (!provider) throw new Error("No contact found for", id);
  globalData.set(id, {
    name: provider.name,
    description: provider.description,
    embeddings: {
      "2D": [...provider.embeddings["2D"], [Math.random(), Math.random()]],
      "3D": [
        ...provider.embeddings["3D"],
        [Math.random(), Math.random(), Math.random()],
      ],
    },
  });
  return new Map(globalData);
}

// fake a cache so we don't slow down stuff we've already seen
let fakeCache = {};

async function fakeNetwork(key) {
  if (!key) {
    fakeCache = {};
  }

  if (fakeCache[key]) {
    return;
  }

  fakeCache[key] = true;
  return new Promise((res) => {
    setTimeout(res, Math.random() * 800);
  });
}
