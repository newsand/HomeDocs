# ID-Center — relicário

Fonte única do visual. A tese aprovada vence qualquer paleta externa.

Descartado da busca de estilo: minimalismo suíço, fundo claro, Inter, Space Grotesk, hero de landing, hover de 200–250ms. Mantido só o que não muda a tese: foco visível, `cursor: pointer`, contraste medido, `prefers-reduced-motion`, sem emoji como ícone.

Não há modo claro.

## Tese visual

Vitrine escura: fundo `#10080C`, painel `#2A0C16`, vinho `#4A1528`, fio de ouro `#C6A15B` e texto em prata `#D5D6D4`. Título em Cormorant Garamond, corpo em Source Serif 4. A página guarda 48px em volta da caixa; dentro, o ritmo é 8 / 12 / 16. Canto de 2px, borda de 1px, sem sombra e sem pílula.

## Tese de movimento

Transição seca de 150ms, `cubic-bezier(0.2, 0, 0, 1)`. O hover só troca a borda de `#C6A15B` para `#E4C98A`. Sem scroll, sem escala, sem bounce, sem elástico. Nada anima `width`, `height`, `top`, `left`, `margin` ou `padding`.

## Cor

| token | hex | uso | contraste |
|---|---|---|---|
| `--ground` | `#10080C` | fundo da página | — |
| `--panel` | `#2A0C16` | caixa | — |
| `--wine` | `#4A1528` | vinho, selo | — |
| `--gold` | `#C6A15B` | borda, fio | 8.14:1 em `--ground`, 7.44:1 em `--panel` |
| `--gold-hover` | `#E4C98A` | borda no hover, anel de foco | 12.26:1 em `--ground` |
| `--silver` | `#D5D6D4` | texto | 13.56:1 em `--ground`, 12.38:1 em `--panel`, 10.08:1 em `--wine` |
| `--silver-dim` | `#9A9792` | rótulo | 6.79:1 em `--ground`, 6.21:1 em `--panel` |

Erro de formulário usa a frase em `--silver` dentro do painel. Sem verde, amarelo ou vermelho novo.

## Tipo

| token | valor |
|---|---|
| `--font-display` | `"Cormorant Garamond", Georgia, serif` |
| `--font-body` | `"Source Serif 4", Georgia, serif` |
| `--size-display` | 48px / 1.05 / 500 |
| `--size-title` | 32px / 1.15 / 500 |
| `--size-body` | 16px / 1.45 / 400 |
| `--size-label` | 13px / 1.45 / 400 |

https://fonts.googleapis.com/css2?family=Cormorant+Garamond:ital,wght@0,500;0,600;1,500&family=Source+Serif+4:opsz,wght@8..60,400;8..60,600&display=swap

## Espaço

| token | valor |
|---|---|
| `--space-8` | 8px |
| `--space-12` | 12px |
| `--space-16` | 16px |
| `--space-24` | 24px |
| `--space-32` | 32px |
| `--space-48` | 48px |

## Canto e sombra

| token | valor |
|---|---|
| `--radius` | 2px |
| `--line` | 1px |
| `--shadow` | none |

## Movimento

| token | valor |
|---|---|
| `--duration-fast` | 150ms |
| `--ease-out` | cubic-bezier(0.2, 0, 0, 1) |

`prefers-reduced-motion: reduce` zera a transição. O hover ainda troca a cor da borda, sem animar.

## Componentes

Borda `--line` `--gold`. Fundo transparente no botão e no campo. Caixa em `--panel`. Selo em `--wine`.

| estado | botão, campo, link |
|---|---|
| default | borda `--gold`, texto `--silver` |
| hover | borda `--gold-hover`, sem escala |
| focus | anel `--gold-hover`, afastado `--radius` |
| active | borda `--silver` |
| disabled | texto e borda `--silver-dim`, sem clique |

O link usa `border-bottom`, não sublinhado de outra cor. O selo não tem hover de escala.

Foco nunca remove o anel sem pôr outro. Todo controle clicável usa `cursor: pointer`.
