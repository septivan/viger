const genrePalette: Record<string, [string, string]> = {
  RPG: ["#8b5cf6", "#ec4899"],
  Roguelike: ["#ef4444", "#f59e0b"],
  Platformer: ["#06b6d4", "#3b82f6"],
  Simulation: ["#10b981", "#84cc16"],
  Adventure: ["#f59e0b", "#f97316"],
  "Action RPG": ["#dc2626", "#7c3aed"],
  Metroidvania: ["#6366f1", "#0ea5e9"],
  Puzzle: ["#14b8a6", "#8b5cf6"],
  Strategy: ["#eab308", "#ef4444"],
  Action: ["#f43f5e", "#8b5cf6"],
  Survival: ["#0f766e", "#65a30d"],
  Sandbox: ["#ca8a04", "#16a34a"],
  Shooter: ["#ea580c", "#dc2626"],
  Racing: ["#2563eb", "#db2777"],
};

export function GameCover({ title, genre, large = false }: { title: string; genre: string; large?: boolean }) {
  const palette = genrePalette[genre] ?? ["#7c3aed", "#2563eb"];
  const initials = title.split(/\s+/).filter((word) => !["the", "of", "and"].includes(word.toLowerCase())).slice(0, 2).map((word) => word[0]).join("");
  return (
    <div className={`game-cover${large ? " game-cover-large" : ""}`} style={{ "--cover-a": palette[0], "--cover-b": palette[1] } as React.CSSProperties}>
      <span className="cover-grid" aria-hidden="true" />
      <span className="cover-orb" aria-hidden="true" />
      <span className="cover-genre">{genre}</span>
      <strong aria-hidden="true">{initials}</strong>
      <span className="cover-title">{title}</span>
    </div>
  );
}

