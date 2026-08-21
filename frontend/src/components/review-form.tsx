"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { APIError, createReview } from "@/lib/api";
import { queryKeys } from "@/lib/query";
import { reviewFormSchema, type ReviewFormValues } from "@/lib/schemas";

export function ReviewForm({ gameID, gameTitle }: { gameID: string; gameTitle: string }) {
  const queryClient = useQueryClient();
  const [submitted, setSubmitted] = useState(false);
  const form = useForm<ReviewFormValues>({
    resolver: zodResolver(reviewFormSchema),
    defaultValues: { reviewerName: "", rating: 0, text: "" },
  });
  const mutation = useMutation({
    mutationFn: (values: ReviewFormValues) => createReview(gameID, values),
    onSuccess: async () => {
      form.reset();
      setSubmitted(true);
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: queryKeys.game(gameID) }),
        queryClient.invalidateQueries({ queryKey: queryKeys.reviewPrefix(gameID) }),
      ]);
    },
    onError: (error) => {
      if (error instanceof APIError) {
        for (const [field, message] of Object.entries(error.fields)) {
          if (field === "reviewerName" || field === "rating" || field === "text") form.setError(field, { message });
        }
      }
    },
  });

  return (
    <section className="review-form-panel" aria-labelledby="write-review-title">
      <p className="eyebrow">Your turn</p>
      <h2 id="write-review-title">What did you think?</h2>
      <p>Share a clear, useful take on <strong>{gameTitle}</strong>.</p>
      <form onSubmit={form.handleSubmit((values) => { setSubmitted(false); mutation.mutate(values); })} noValidate>
        <fieldset>
          <legend>Your rating</legend>
          <div className="star-input">
            <Controller control={form.control} name="rating" render={({ field }) => <>{[1, 2, 3, 4, 5].map((rating) => <label key={rating}><input checked={field.value === rating} name={field.name} onBlur={field.onBlur} onChange={() => field.onChange(rating)} ref={field.ref} type="radio" value={rating} /><span aria-hidden="true">★</span><span className="sr-only">{rating} out of 5</span></label>)}</>} />
          </div>
          {form.formState.errors.rating && <p className="field-error">{form.formState.errors.rating.message}</p>}
        </fieldset>
        <label className="form-field"><span>Your name</span><input autoComplete="name" maxLength={80} placeholder="How should we credit you?" {...form.register("reviewerName")} />{form.formState.errors.reviewerName && <small className="field-error">{form.formState.errors.reviewerName.message}</small>}</label>
        <label className="form-field"><span>Your review</span><textarea maxLength={2000} placeholder="What worked, what did not, and who would enjoy it?" rows={6} {...form.register("text")} />{form.formState.errors.text && <small className="field-error">{form.formState.errors.text.message}</small>}</label>
        {mutation.isError && <p className="form-message error" role="alert">{mutation.error instanceof Error ? mutation.error.message : "Your review could not be submitted."}</p>}
        {submitted && <p className="form-message success" role="status">Your review is now part of the conversation.</p>}
        <button className="primary-button" disabled={mutation.isPending} type="submit">{mutation.isPending ? "Publishing…" : "Publish review"}<span aria-hidden="true">↗</span></button>
      </form>
    </section>
  );
}
