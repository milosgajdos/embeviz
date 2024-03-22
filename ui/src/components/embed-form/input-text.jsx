export default function InputText({ text, onTextChange }) {
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
