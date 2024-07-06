import { useEffect } from "react";
import type { PortfolioItem } from "@components/Headless/Portfolio/PortfolioItems";

interface DialogPopupProps {
  open: boolean;
  item: PortfolioItem | null;
  onClose: () => void;
}
export default function DialogPopup({ open, item, onClose }: DialogPopupProps) {
  useEffect(() => {
    const handleEsc = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };

    if (open) {
      document.addEventListener("keydown", handleEsc);
    } else {
      document.removeEventListener("keydown", handleEsc);
    }

    return () => document.removeEventListener("keydown", handleEsc);
  }, [open, onClose]);

  if (!open || !item) return null;

  return (
    <dialog className="dialog fixed z-50 cursor-pointer" open={open} onClick={onClose}>
      <div className="flex justify-center items-center h-screen">
        <img className="image" src={item.src} alt={item.alt} />
      </div>
    </dialog>
  )
}
