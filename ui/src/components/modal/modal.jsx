import { useState, useRef, useEffect } from "react";

// isOpen controls if the modal is shown or not
// onClose allows the parent to pass down onClose handler
// children allows the parent to pass down children elements e.g. buttons, etc.
export default function Modal({ isOpen, onClose, children }) {
  const [isModalOpen, setModalOpen] = useState(isOpen);
  const modalRef = useRef(null);

  useEffect(() => {
    setModalOpen(isOpen);
  }, [isOpen]);

  useEffect(() => {
    const modalElement = modalRef.current;
    if (modalElement) {
      if (isModalOpen) {
        modalElement.showModal();
      } else {
        modalElement.close();
      }
    }
  }, [isModalOpen]);

  const handleCloseModal = () => {
    if (onClose) {
      onClose();
    }
    setModalOpen(false);
  };

  const handleKeyDown = (e) => {
    if (e.key === "Escape") {
      handleCloseModal();
    }
  };

  const handleOutsideClick = (e) => {
    if (modalRef.current && isModalOpen) {
      handleCloseModal();
    }
  };

  return (
    <div className="modal-container" onClick={handleOutsideClick}>
      <dialog ref={modalRef} onKeyDown={handleKeyDown} className="modal">
        {children}
      </dialog>
    </div>
  );
}