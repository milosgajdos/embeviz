import { matchSorter } from "match-sorter";
import sortBy from "sort-by";

const defaultChunking = {
  size: 2,
  overlap: 0,
  trim: false,
  sep: false,
};

export async function getProviders(query) {
  const resp = await fetch("http://localhost:5050/api/v1/providers");
  const respData = await resp.json();
  let providers = respData.providers;
  if (!providers) return [];
  if (query) {
    providers = matchSorter(providers, query, { keys: ["name"] });
  }
  return providers.sort(sortBy("name"));
}

export async function getProvider(uid) {
  const resp = await fetch("http://localhost:5050/api/v1/providers/" + uid);
  const provider = await resp.json();
  return provider ?? null;
}

export async function getProviderProjections(uid) {
  const resp = await fetch(
    "http://localhost:5050/api/v1/providers/" + uid + "/projections",
  );
  const embeddings = await resp.json();
  return embeddings ?? null;
}

export async function embedData(uid, updates) {
  let data = {
    text: updates.text,
    label: updates.label,
    projection: updates.projection,
  };

  if (updates.chunking === "on") {
    data.chunking = defaultChunking;

    if (updates.size) {
      data.chunking.size = parseInt(updates.size, 10);
    }
    if (updates.chunking.overlap) {
      data.chunking.overlap = parseInt(updates.overlap, 10);
    }

    if (updates.trim === "on") {
      data.chunking.trim = true;
    }
    if (updates.sep === "on") {
      data.chunking.sep = true;
    }
  }

  await fetch("http://localhost:5050/api/v1/providers/" + uid + "/embeddings", {
    method: "PUT",
    body: JSON.stringify(data),
    headers: {
      "Content-Type": "application/json",
    },
  });
}

export async function deleteData(uid) {
  await fetch("http://localhost:5050/api/v1/providers/" + uid + "/embeddings", {
    method: "DELETE",
  });
}

export async function computeData(uid, updates) {
  await fetch(
    "http://localhost:5050/api/v1/providers/" + uid + "/projections",
    {
      method: "PATCH",
      body: JSON.stringify(updates),
      headers: {
        "Content-Type": "application/json",
      },
    },
  );
}
