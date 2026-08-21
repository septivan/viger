import { z } from "zod";

export const paginationSchema = z.object({
  page: z.number().int().positive(),
  pageSize: z.number().int().positive(),
  totalItems: z.number().int().nonnegative(),
  totalPages: z.number().int().nonnegative(),
});

export const gameSchema = z.object({
  id: z.string(),
  title: z.string(),
  description: z.string(),
  genre: z.string(),
  platforms: z.array(z.string()),
  developer: z.string(),
  releaseDate: z.string(),
  averageRating: z.number().min(0).max(5),
  reviewCount: z.number().int().nonnegative(),
  ratingDistribution: z.record(z.string(), z.number().int().nonnegative()),
});

export const reviewSchema = z.object({
  id: z.string(),
  gameId: z.string(),
  reviewerName: z.string(),
  rating: z.number().int().min(1).max(5),
  text: z.string(),
  createdAt: z.string(),
});

export const gamePageSchema = z.object({ data: z.array(gameSchema), pagination: paginationSchema });
export const gameResponseSchema = z.object({ data: gameSchema });
export const reviewPageSchema = z.object({ data: z.array(reviewSchema), pagination: paginationSchema });
export const reviewResponseSchema = z.object({ data: reviewSchema });

export const reviewFormSchema = z.object({
  reviewerName: z.string().trim().min(2, "Enter at least 2 characters.").max(80, "Use no more than 80 characters."),
  rating: z.number().int().min(1, "Choose a rating.").max(5),
  text: z.string().trim().min(10, "Enter at least 10 characters.").max(2000, "Use no more than 2,000 characters."),
});

export const reviewCreatedEventSchema = z.object({
  type: z.literal("review.created"),
  gameId: z.string(),
  review: reviewSchema,
});

export type Game = z.infer<typeof gameSchema>;
export type Review = z.infer<typeof reviewSchema>;
export type Pagination = z.infer<typeof paginationSchema>;
export type ReviewFormValues = z.infer<typeof reviewFormSchema>;
export type ReviewCreatedEvent = z.infer<typeof reviewCreatedEventSchema>;

