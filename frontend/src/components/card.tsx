"use client"

import Image from "next/image"
import { useState } from "react"

export type CardData = {
  cardImageURL: string
  cardRarity: "common" | "uncommon" | "rare" | "mythic"
}

export type CardImageProps = CardData

export function CardImage(props: CardImageProps) {
  let [isFaceUp, setIsFaceUp] = useState<boolean>(false)
  let borderRarityColor = {
    "common": "border-rarity-common",
    "uncommon": "border-rarity-uncommon",
    "rare": "border-rarity-rare",
    "mythic": "border-rarity-mythic",
  }[props.cardRarity]
  let shadowRarityColor = {
    "common": "shadow-rarity-common",
    "uncommon": "shadow-rarity-uncommon",
    "rare": "shadow-rarity-rare",
    "mythic": "shadow-rarity-mythic",
  }[props.cardRarity]
  return (
    <div className="group w-52 h-72 hover:scale-150 duration-100 z-[100] hover:z-[110] [perspective:1000px]">
      <div onClick={() => { if (!isFaceUp) setIsFaceUp(true) }} className={
        `absolute rounded-xl h-full w-full duration-500 transition-all [transform-style:preserve-3d] ${!isFaceUp && "[transform:rotateY(180deg)]"}`
      }>
        <div className="absolute inset-0 [backface-visibility:hidden]">
          <Image className={`duration-100 border-2 ${borderRarityColor} w-full h-full rounded-xl ${isFaceUp && "shadow-[0px_0px_20px_1px_rgba(0,0,0,0.3)]"} ${shadowRarityColor}`} src={props.cardImageURL} alt="" width={1024} height={1024} quality={100} layout="" />
        </div>
        <div className="absolute inset-0 [backface-visibility:hidden] [transform:rotateY(180deg)]">
          <Image className="duration-100 border-2 border-white rounded-xl" src="/cardback.webp" alt="" width={1024} height={1024} quality={100} layout="" />
        </div>
      </div>
    </div >
  )
}

export type CardDisplayProps = {
  cardImageURLs: CardData[]
}

export function CardDisplay(props: CardDisplayProps) {
  return (
    <div className="flex flex-wrap flex-row mx-16 my-16 gap-6 items-center justify-center">
      {props.cardImageURLs.map((cardData: CardData, i: number) => {
        return (
          <CardImage key={i} {...cardData} />
        )
      })}
    </div>
  )
}