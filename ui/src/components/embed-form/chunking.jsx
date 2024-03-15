export default function Chunking({
  chunking,
  onChunkingChange,
  size,
  onSizeChange,
  overlap,
  onOverlapChange,
}) {
  return (
    <div id="embed-chunking">
      <fieldset>
        <legend>Chunking</legend>
        <input
          type="checkbox"
          id="chunking"
          name="chunking"
          checked={chunking}
          onChange={onChunkingChange}
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
                onChange={onSizeChange}
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
                onChange={onOverlapChange}
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
  );
}
