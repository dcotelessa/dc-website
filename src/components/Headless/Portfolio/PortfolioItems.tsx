import { useState } from "react";
import { motion } from "framer-motion";
import DialogPopup from "@components/Portfolio/DialogPopup";

export interface PortfolioItem {
  src: string;
  alt?: string;
  width?: number;
  height?: number;
}

interface PortfolioItemsProps {
  items: PortfolioItem[];
}

export function PortfolioItems({ items }: PortfolioItemsProps) {
  const [selectedItem, setSelectedItem] = useState<PortfolioItem | null>(null);
  const [loadedImages, setLoadedImages] = useState<{ [key: string]: boolean }>({});

  const openDialog = (item: PortfolioItem) => {
    setSelectedItem(item);
  };

  const closeDialog = () => {
    setSelectedItem(null);
  };

  const handleImageLoad = (alt: string) => {
    setLoadedImages((prev) => ({ ...prev, [alt]: true }));
  };

  const animationVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0 },
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {items.map((item) => (
        <div key={item.alt} className="portfolio-item">
          <a onClick={() => openDialog(item)} className="cursor-pointer">
            <motion.img 
              src={item.src} 
              alt={item.alt} 
              width={item.width} 
              height={item.height}
              loading="lazy"
              onLoad={() => handleImageLoad(item.alt || "")}
              initial="hidden"
              animate={loadedImages[item.alt || ""] ? "visible" : "hidden"}
              variants={animationVariants}
              transition={{ duration: 0.6, ease: "easeOut" }}
              whileHover={{ scale: 1.1, transition: { duration: 0.2 } }} // This replaces the CSS hover effect
            />
          </a>
        </div>
      ))}
      <DialogPopup open={!!selectedItem} item={selectedItem} onClose={closeDialog} />
    </div>
  );
}

export default PortfolioItems;
