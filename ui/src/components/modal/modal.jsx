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

  // TODO: this does not seem to be working
  // I've also tried !(e.target.contains(modalRef.current)'
  // it seems e.target keeps getting captured by dialog instead of modal-container
  const handleOutsideClick = (e) => {
    console.log(e.target.className);
    if (
      modalRef.current &&
      isModalOpen &&
      e.target.className === "modal-container"
    ) {
      handleCloseModal();
    }
  };

  return (
    <div className="modal-container" onClick={handleOutsideClick}>
      <dialog className="modal" ref={modalRef} onKeyDown={handleKeyDown}>
        {children}
      </dialog>
    </div>
  );
}
