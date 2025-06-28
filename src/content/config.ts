import { defineCollection, z } from 'astro:content';

const solutionsCollection = defineCollection({
  schema: z.object({
    title: z.string(),
    target: z.number(),
    date: z.date(),
    bytes: z.number(),
    difficulty: z.enum(['easy', 'medium', 'hard']),
    month: z.string(), // "2024-01" for folder grouping
  }),
});

export const collections = {
  solutions: solutionsCollection,
};
