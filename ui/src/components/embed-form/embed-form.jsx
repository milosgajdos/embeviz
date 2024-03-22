import { Form, useNavigation, useParams } from "react-router-dom";
import { useState, useEffect } from "react";
import { deleteData } from "../../lib/embeddings";
import InputLabel from "./input-label";
import InputText from "./input-text";
import Projection from "./projection";
import Chunking from "./chunking";
import Modal from "../modal/modal";

export default function EmbedForm({ onDrop, onFetch }) {
  let params = useParams();
  const navigation = useNavigation();

  // Input fields
  const [label, setLabel] = useState("");
  const [text, setText] = useState("");
  // Projection options
  const [projection, setProjection] = useState("pca");
  const [color, setColor] = useState("#0000FF");
  // Chunking options
  const [chunking, setChunking] = useState(false);
  const [size, setSize] = useState("2");
  const [overlap, setOverlap] = useState("0");
  // Drop data modal
  const [isOpenModal, setOpenModal] = useState(false);

  useEffect(() => {
    onFetch(
      navigation.state === "submitting" || navigation.state === "loading",
    );
  }, [navigation.state]);

  async function handleDrop() {
    try {
      await deleteData(params.uid);
      onDrop();
      setOpenModal(false);
    } catch (error) {
      console.error("Error deleting data:", error);
    }
  }

  return (
    <>
      <Form method="post" id="embed-form">
        <div id="embed-input">
          <InputLabel
            label={label}
            onLabelChange={(e) => setLabel(e.target.value)}
          />
          <InputText
            text={text}
            onTextChange={(e) => setText(e.target.value)}
          />
        </div>
        <Projection
          projection={projection}
          onProjectionChange={(e) => setProjection(e.target.value)}
          onProjectionClearInput={() => {
            setLabel("");
            setText("");
          }}
          onProjetionDrop={() => setOpenModal(true)}
          color={color}
          onColorChange={(e) => setColor(e.target.value)}
        />
        <Chunking
          chunking={chunking}
          onChunkingChange={() => setChunking(!chunking)}
          size={size}
          onSizeChange={(e) => setSize(e.target.value)}
          overlap={overlap}
          onOverlapChange={(e) => setOverlap(e.target.value)}
        />
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
              onClick={handleDrop}
            >
              OK
            </button>
          </div>
        </div>
      </Modal>
    </>
  );
}
