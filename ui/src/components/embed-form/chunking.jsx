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
        <div>
          <input
            type="checkbox"
            id="chunking"
            name="chunking"
            checked={chunking}
            onChange={onChunkingChange}
          />
          <label htmlFor="chunking"> Enable</label>
        </div>
        <fieldset disabled={!chunking}>
          <legend>Options</legend>
          <div className="chunking-splits">
            <div className="chunking-size">
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
            <div className="chunking-overlap">
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
          </div>
          <div>
            <input type="checkbox" id="trim" name="trim" />
            <label htmlFor="trim"> Trim</label>
          </div>
          <div>
            <input type="checkbox" id="sep" name="sep" />
            <label htmlFor="sep"> Separator</label>
          </div>
          <button
            type="submit"
            className="update-btn"
            name="intent"
            value="chunk"
          >
            Highlight
          </button>
        </fieldset>
      </fieldset>
    </div>
  );
}
