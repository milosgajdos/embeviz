export default function InputLabel({ label, onLabelChange }) {
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
