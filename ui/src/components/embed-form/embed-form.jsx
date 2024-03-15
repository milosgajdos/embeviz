import { Form, useNavigation, useParams } from "react-router-dom";
import { useState, useEffect } from "react";
import { deleteData } from "../../lib/embeddings";
import Modal from "../modal/modal";

export default function EmbedForm({ onDeletion, onFetching }) {
  let params = useParams();
  const navigation = useNavigation();
  const [projection, setProjection] = useState("pca");
  const [chunking, setChunking] = useState(false);
  const [size, setSize] = useState("2");
  const [overlap, setOverlap] = useState("0");
  const [label, setLabel] = useState("");
  const [text, setText] = useState("");
  const [isOpenModal, setOpenModal] = useState(false);

  useEffect(() => {
    onFetching(
      navigation.state === "submitting" || navigation.state === "loading",
    );
  }, [navigation.state]);

  async function handleDeletion() {
    try {
      await deleteData(params.uid);
      onDeletion();
      setOpenModal(false);
    } catch (error) {
      console.error("Error deleting data:", error);
    }
  }

  function handleProjection(e) {
    setProjection(e.target.value);
  }

  function handleClearFields() {
    setLabel("");
    setText("");
  }

  return (
    <>
      <Form method="post" id="embed-form">
        <div id="embed-input">
          <input
            id="label"
            name="label"
            placeholder="Label (Optional)"
            value={label}
            onChange={(e) => setLabel(e.target.value)}
          />
          <textarea
            id="text"
            name="text"
            placeholder="Text"
            rows="10"
            cols="80"
            wrap="soft"
            value={text}
            onChange={(e) => setText(e.target.value)}
          ></textarea>
        </div>
        <div id="embed-projection">
          <fieldset>
            <legend>Projection</legend>
            <div>
              <input
                type="radio"
                id="pca"
                name="projection"
                value="pca"
                checked={projection === "pca"}
                onChange={handleProjection}
              />
              <label htmlFor="pca"> pca</label>
            </div>
            <div>
              <input
                type="radio"
                id="tsne"
                name="projection"
                value="tsne"
                checked={projection === "tsne"}
                onChange={handleProjection}
              />
              <label htmlFor="tsne"> t-SNE</label>
            </div>
          </fieldset>
          <button type="submit" name="intent" value="embed">
            Embed
          </button>
          <button
            type="submit"
            className="update-btn"
            name="intent"
            value="compute"
          >
            Recompute
          </button>
          <button
            type="button"
            className="delete-btn"
            onClick={handleClearFields}
          >
            Clear
          </button>
          <button
            type="button"
            className="delete-btn"
            name="modal"
            value="modal"
            onClick={() => setOpenModal(true)}
          >
            Drop
          </button>
        </div>
        <div id="embed-chunking">
          <fieldset>
            <legend>Chunking</legend>
            <input
              type="checkbox"
              id="chunking"
              name="chunking"
              checked={chunking}
              onChange={() => setChunking(!chunking)}
            />
            <label htmlFor="chunking"> Enable</label>
            <fieldset disabled={!chunking}>
              <legend>Options</legend>
              <div>
                <div>
                  <label htmlFor="size">Size </label>
                  <input
                    type="number"
                    id="size"
                    name="size"
                    min="2"
                    value={size}
                    onChange={(e) => setSize(e.target.value)}
                  />
                </div>
                <div>
                  <label htmlFor="overlap">Overlap </label>
                  <input
                    type="number"
                    id="overlap"
                    name="overlap"
                    min="0"
                    value={overlap}
                    onChange={(e) => setOverlap(e.target.value)}
                  />
                </div>
                <br />
                <div>
                  <input type="checkbox" id="trim" name="trim" />
                  <label htmlFor="trim"> Trim</label>
                </div>
                <div>
                  <input type="checkbox" id="sep" name="sep" />
                  <label htmlFor="sep"> Separator</label>
                </div>
              </div>
            </fieldset>
          </fieldset>
        </div>
      </Form>
      <Modal isOpen={isOpenModal} onClose={() => setOpenModal(false)}>
        <div className="modal-content">
          <p> Are you sure you want to drop the data? </p>
          <div className="modal-buttons">
            <button
              value="cancel"
              className="modal-cancel-btn"
              onClick={() => setOpenModal(false)}
            >
              Cancel
            </button>
            <button
              value="ok"
              className="modal-ok-btn delete-btn"
              onClick={handleDeletion}
            >
              OK
            </button>
          </div>
        </div>
      </Modal>
    </>
  );
}
