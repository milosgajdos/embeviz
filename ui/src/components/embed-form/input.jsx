import { HighlightWithinTextarea } from "react-highlight-within-textarea";

export default function Input({
  label,
  onLabelChange,
  text,
  onTextChange,
  chunks,
  hlText,
  onHlText,
}) {
  return (
    <div id="embed-input">
      <InputLabel label={label} onLabelChange={onLabelChange} />
      <InputText
        text={text}
        onTextChange={onTextChange}
        chunks={chunks}
        hlText={hlText}
        onHlText={onHlText}
      />
    </div>
  );
}

function InputLabel({ label, onLabelChange }) {
  return (
    <input
      id="label"
      name="label"
      placeholder="Label (Optional)"
      value={label}
      onChange={onLabelChange}
    />
  );
}

function InputText({ text, onTextChange, chunks, hlText, onHlText }) {
  return (
    <>
      {/* <textarea
        id="text"
        name="text"
        placeholder="Input"
        rows="10"
        cols="80"
        wrap="soft"
        value={text}
        onChange={onTextChange}
      ></textarea> */}
      <div className="area" id="text">
        <HighlightWithinTextarea
          name="text"
          value={hlText}
          onChange={onHlText}
          highlight={chunks}
        />
      </div>
    </>
  );
}
