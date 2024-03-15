export default function Input({ label, onLabelChange, text, onTextChange }) {
  return (
    <div id="embed-input">
      <input
        id="label"
        name="label"
        placeholder="Label (Optional)"
        value={label}
        onChange={onLabelChange}
      />
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
    </div>
  );
}
