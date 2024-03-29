export default function Projection({
  projection,
  onProjectionChange,
  onProjectionClearInput,
  onProjetionDrop,
  color,
  onColorChange,
}) {
  return (
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
            onChange={onProjectionChange}
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
            onChange={onProjectionChange}
          />
          <label htmlFor="tsne"> t-SNE</label>
        </div>
        <div>
          <input
            type="color"
            id="color"
            name="color"
            value={color}
            onChange={onColorChange}
          />
          <label htmlFor="color"> Color</label>
        </div>
      </fieldset>
      <button type="submit" className="embed-btn" name="intent" value="embed">
        Embed
      </button>
      <button
        type="submit"
        className="update-btn"
        name="intent"
        value="compute"
      >
        Compute
      </button>
      <button
        type="button"
        className="delete-btn"
        name="clear"
        value="clear"
        onClick={onProjectionClearInput}
      >
        Clear
      </button>
      <button
        type="button"
        className="delete-btn"
        name="modal"
        value="modal"
        onClick={onProjetionDrop}
      >
        Drop
      </button>
    </div>
  );
}
