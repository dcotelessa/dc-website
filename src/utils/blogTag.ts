export const tags = ["coding", "technology", "nerd"] as const;

export type TagItem = {
  href: string;
  color: string;
};

export type ValidTags = (typeof tags)[number];

type TagMap = {
  [key in ValidTags]: TagItem;
};

export const tagMap: TagMap = {
  coding: {
    href: "/blog/tags/coding",
    color: "blue-400",
  },
  technology: {
    href: "/blog/tags/technology",
    color: "teal-400",
  },
  nerd: {
    href: "/blog/tags/nerd",
    color: "cyan-200",
  },
};
