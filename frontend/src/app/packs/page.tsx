import { CardDisplay, CardImage } from "@/components/card";
import Layout from "@/components/layout";

export default function PacksPage() {
  return (
    <Layout>
      <CardDisplay cardImageURLs={[
        {
          cardImageURL: "https://cards.scryfall.io/png/front/3/4/343d01cf-9806-4c2d-a993-ddc9ed248d7f.png",
          cardRarity: "rare"
        },
        {
          cardImageURL: "https://cards.scryfall.io/large/front/4/3/434515bf-de57-4c00-b0b4-c9579cc1b84c.jpg",
          cardRarity: "uncommon"
        },
        {
          cardImageURL: "https://cards.scryfall.io/large/front/4/3/434515bf-de57-4c00-b0b4-c9579cc1b84c.jpg",
          cardRarity: "uncommon"
        },
        {
          cardImageURL: "https://cards.scryfall.io/png/front/9/a/9afe8b9e-bb14-44d5-b5da-627835ee457f.png",
          cardRarity: "common"
        },
        {
          cardImageURL: "https://cards.scryfall.io/png/front/9/a/9afe8b9e-bb14-44d5-b5da-627835ee457f.png",
          cardRarity: "common"
        },
        {
          cardImageURL: "https://cards.scryfall.io/png/front/9/a/9afe8b9e-bb14-44d5-b5da-627835ee457f.png",
          cardRarity: "common"
        },
        {
          cardImageURL: "https://cards.scryfall.io/png/front/9/a/9afe8b9e-bb14-44d5-b5da-627835ee457f.png",
          cardRarity: "common"
        },

      ]} />
    </Layout>
  )
}