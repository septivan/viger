import Link from "next/link";

export default function NotFound() {
  return <main className="detail-page"><div className="state-panel"><span>404</span><h1>This path leads nowhere.</h1><p>The game or page you were looking for is not in the collection.</p><Link href="/">Return to Viger</Link></div></main>;
}

