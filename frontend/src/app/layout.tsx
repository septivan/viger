import type { Metadata } from "next";
import Link from "next/link";
import { Providers } from "@/app/providers";
import "@/app/styles.css";

export const metadata: Metadata = {
  title: { default: "Viger — Find your next great game", template: "%s — Viger" },
  description: "Browse games, share thoughtful reviews, and see new opinions arrive live.",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <Providers>
          <header className="site-header">
            <Link className="brand" href="/" aria-label="Viger home">
              <span className="brand-mark" aria-hidden="true">V</span>
              <span>Viger</span>
            </Link>
            <p>Reviews from people who play.</p>
          </header>
          {children}
          <footer className="site-footer">
            <span>Viger</span>
            <p>Discover, rate, and share.</p>
          </footer>
        </Providers>
      </body>
    </html>
  );
}
