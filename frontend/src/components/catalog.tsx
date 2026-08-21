"use client";

import { useQuery } from "@tanstack/react-query";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useCallback, useEffect, useMemo, useState } from "react";
import { GameCard } from "@/components/game-card";
import { Pagination } from "@/components/pagination";
import { listGames, type GameQuery } from "@/lib/api";
import { queryKeys } from "@/lib/query";

const genres = ["", "Action", "Action RPG", "Adventure", "Metroidvania", "Platformer", "Puzzle", "Racing", "Roguelike", "RPG", "Sandbox", "Shooter", "Simulation", "Strategy", "Survival"];
const platforms = ["", "PC", "Nintendo Switch", "PlayStation 5", "Xbox Series", "Mobile"];

export function Catalog() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const currentSearch = searchParams.get("q") ?? "";
  const [search, setSearch] = useState(currentSearch);

  const updateParameters = useCallback((updates: Record<string, string | number | undefined>) => {
    const next = new URLSearchParams(searchParams.toString());
    for (const [key, value] of Object.entries(updates)) {
      if (value === undefined || value === "" || value === 0) next.delete(key);
      else next.set(key, String(value));
    }
    router.replace(`${pathname}${next.size ? `?${next}` : ""}`, { scroll: false });
  }, [pathname, router, searchParams]);

  useEffect(() => {
    if (search === currentSearch) return;
    const timeout = window.setTimeout(() => updateParameters({ q: search.trim(), page: 1 }), 300);
    return () => window.clearTimeout(timeout);
  }, [currentSearch, search, updateParameters]);

  const query = useMemo<GameQuery>(() => ({
    q: currentSearch || undefined,
    genre: searchParams.get("genre") || undefined,
    platform: searchParams.get("platform") || undefined,
    minRating: Number(searchParams.get("minRating") ?? 0),
    sort: searchParams.get("sort") || "rating_desc",
    page: Math.max(1, Number(searchParams.get("page") ?? 1) || 1),
    pageSize: 12,
  }), [currentSearch, searchParams]);

  const games = useQuery({ queryKey: queryKeys.games(query), queryFn: () => listGames(query) });

  return (
    <section className="catalog-section" id="catalog">
      <div className="section-heading">
        <div><p className="eyebrow">The collection</p><h2>Worth playing.<br /><em>Worth discussing.</em></h2></div>
        <p>{games.data ? `${games.data.pagination.totalItems} games selected for curious players.` : "A considered collection across genres and platforms."}</p>
      </div>
      <div className="catalog-controls">
        <label className="search-field"><span className="sr-only">Search games</span><span aria-hidden="true">⌕</span><input onChange={(event) => setSearch(event.target.value)} placeholder="Search by title…" type="search" value={search} /></label>
        <label><span>Genre</span><select onChange={(event) => updateParameters({ genre: event.target.value, page: 1 })} value={query.genre ?? ""}>{genres.map((genre) => <option key={genre || "all"} value={genre}>{genre || "All genres"}</option>)}</select></label>
        <label><span>Platform</span><select onChange={(event) => updateParameters({ platform: event.target.value, page: 1 })} value={query.platform ?? ""}>{platforms.map((platform) => <option key={platform || "all"} value={platform}>{platform || "All platforms"}</option>)}</select></label>
        <label><span>Sort by</span><select onChange={(event) => updateParameters({ sort: event.target.value, page: 1 })} value={query.sort}><option value="rating_desc">Top rated</option><option value="reviews_desc">Most discussed</option><option value="newest">Newest releases</option><option value="title_asc">Title A–Z</option></select></label>
      </div>
      {games.isPending && <div className="game-grid" aria-label="Loading games">{Array.from({ length: 8 }, (_, index) => <div className="skeleton game-card-skeleton" key={index} />)}</div>}
      {games.isError && <div className="state-panel" role="alert"><span>Connection lost</span><h3>The collection could not be loaded.</h3><p>{games.error instanceof Error ? games.error.message : "Try again in a moment."}</p><button onClick={() => games.refetch()} type="button">Try again</button></div>}
      {games.data?.data.length === 0 && <div className="state-panel"><span>No matches</span><h3>Try a broader search.</h3><p>Adjust the title, genre, or platform to rediscover the collection.</p><button onClick={() => { setSearch(""); router.replace(pathname, { scroll: false }); }} type="button">Clear filters</button></div>}
      {games.data && games.data.data.length > 0 && <><div className="game-grid">{games.data.data.map((game) => <GameCard game={game} key={game.id} />)}</div><Pagination label="Game pages" onPageChange={(page) => updateParameters({ page })} page={games.data.pagination.page} totalPages={games.data.pagination.totalPages} /></>}
    </section>
  );
}

