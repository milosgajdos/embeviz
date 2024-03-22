export default function Input({ label, onLabelChange, text, onTextChange }) {
  return (
    <div id="embed-input">
      <InputLabel label={label} onLabelChange={onLabelChange} />
      <InputText text={text} onTextChange={onTextChange} />
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

function InputText({ text, onTextChange }) {
  return (
    <textarea
      id="text"
      name="text"
      placeholder="Text"
      rows="10"
      cols="80"
      wrap="soft"
      value={text}
      onChange={onTextChange}
    ></textarea>
  );
}
