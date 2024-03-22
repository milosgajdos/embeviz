import { matchSorter } from "match-sorter";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:5050/api/v1";

const defaultChunking = {
  size: 2,
  overlap: 0,
  trim: false,
  sep: false,
};

export async function getProviders(query) {
  try {
    const resp = await fetch(API_URL + "/providers");
    if (!resp.ok) {
      throw new Error(`HTTP error! Status: ${resp.status}`);
    }
    const respData = await resp.json();
    let providers = respData.providers;
    if (!providers) return [];

    if (!query) {
      query = "";
    }
    return matchSorter(providers, query, { keys: ["name"] });
  } catch (error) {
    console.error("An error occurred:", error.message);
    throw new Error(`Error fetching providers! Message: ${error.message}`);
  }
}

export async function getProvider(uid) {
  try {
    const resp = await fetch(API_URL + "/providers/" + uid);
    if (!resp.ok) {
      throw new Error(`HTTP error! Status: ${resp.status}`);
    }

    const provider = await resp.json();
    return provider ?? null;
  } catch (error) {
    console.error("An error occurred:", error.message);
    throw new Error(
      `Error fetching provider ${uid}! Message: ${error.message}`,
    );
  }
}

export async function getProviderProjections(uid) {
  try {
    const resp = await fetch(API_URL + "/providers/" + uid + "/projections");
    if (!resp.ok) {
      throw new Error(`HTTP error! Status: ${resp.status}`);
    }

    const embeddings = await resp.json();
    return embeddings ?? null;
  } catch (error) {
    console.error("An error occurred:", error.message);
    throw new Error(
      `Error fetching provider ${uid} projections! Message: ${error.message}`,
    );
  }
}

export async function embedData(uid, updates) {
  let data = {
    text: updates.text,
    label: updates.label,
    projection: updates.projection,
    metadata: {
      color: updates.color,
    },
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

  try {
    const resp = await fetch(API_URL + "/providers/" + uid + "/embeddings", {
      method: "PUT",
      body: JSON.stringify(data),
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!resp.ok) {
      throw new Error(`HTTP error! Status: ${resp.status}`);
    }
  } catch (error) {
    console.error("An error occurred:", error.message);
    throw new Error(
      `Error embdding data for provider ${uid}! Message: ${error.message}`,
    );
  }
}

export async function deleteData(uid) {
  try {
    const resp = await fetch(API_URL + "/providers/" + uid + "/embeddings", {
      method: "DELETE",
    });

    if (!resp.ok) {
      throw new Error(`HTTP error! Status: ${resp.status}`);
    }
  } catch (error) {
    console.error("An error occurred:", error.message);
    throw new Error(
      `Error deleting data for provider ${uid}! Message: ${error.message}`,
    );
  }
}

export async function computeData(uid, updates) {
  try {
    const resp = await fetch(API_URL + "/providers/" + uid + "/projections", {
      method: "PATCH",
      body: JSON.stringify(updates),
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!resp.ok) {
      throw new Error(`HTTP error! Status: ${resp.status}`);
    }
  } catch (error) {
    console.error("An error occurred:", error.message);
    throw new Error(
      `Error computing data for provider ${uid}! Message: ${error.message}`,
    );
  }
}
