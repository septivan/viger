import { Suspense } from "react";
import { Catalog } from "@/components/catalog";

export default function HomePage() {
  return (
    <main>
      <section className="hero">
        <div className="hero-orbit" aria-hidden="true"><span /><span /><span /></div>
        <p className="eyebrow"><span className="live-dot" /> Community powered</p>
        <h1>Find the game<br />that stays with you.</h1>
        <p className="hero-copy">Explore honest player perspectives, discover something unexpected, and add your voice while the conversation is live.</p>
        <a className="hero-action" href="#catalog">Browse the collection <span aria-hidden="true">↓</span></a>
      </section>
      <Suspense fallback={<CatalogSkeleton />}>
        <Catalog />
      </Suspense>
    </main>
  );
}

function CatalogSkeleton() {
  return <section className="catalog-section" aria-label="Loading games"><div className="skeleton skeleton-heading" /><div className="game-grid">{Array.from({ length: 8 }, (_, index) => <div className="skeleton game-card-skeleton" key={index} />)}</div></section>;
}

