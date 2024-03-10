"use client"

import { CardData } from "@/types/card"
import Image from "next/image"
import { useState } from "react"

export type CardImageProps = CardData & {
  startFaceUp?: boolean
}

export function CardImage(props: CardImageProps) {
  let [isFaceUp, setIsFaceUp] = useState<boolean>(props.startFaceUp || false)
  let borderRarityColor = {
    "common": "border-rarity-common",
    "uncommon": "border-rarity-uncommon",
    "rare": "border-rarity-rare",
    "mythic": "border-rarity-mythic",
  }[props.card_rarity]
  let shadowRarityColor = {
    "common": "shadow-rarity-common",
    "uncommon": "shadow-rarity-uncommon",
    "rare": "shadow-rarity-rare",
    "mythic": "shadow-rarity-mythic",
  }[props.card_rarity]
  return (
    <div className="group w-[256px] h-[355px] hover:scale-125 scale-100 duration-75 z-[100] hover:z-[110] [perspective:1000px]">
      <div onClick={() => { if (!isFaceUp) setIsFaceUp(true) }} className={
        `absolute rounded-xl w-full h-full duration-500 transition-all [transform-style:preserve-3d] ${!isFaceUp && "[transform:rotateY(180deg)]"}`
      }>
        <div className="absolute inset-0 [backface-visibility:hidden]">
          <Image className={`duration-75 border-2 ${borderRarityColor} rounded-xl ${isFaceUp && "shadow-[0px_0px_20px_1px_rgba(0,0,0,0.3)]"} ${shadowRarityColor}`} unoptimized priority src={props.image_url} alt="" width={1024} height={768} quality={100} />
        </div>
        <div className="absolute inset-0 [backface-visibility:hidden] [transform:rotateY(180deg)]">
          <Image className="duration-75 border-2 border-white rounded-xl" src="/cardback.webp" alt="" width={1024} height={1024} quality={100} layout="" />
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
    <div className="flex flex-wrap flex-row gap-2 items-center justify-center">
      {props.cardImageURLs.map((cardData: CardData, i: number) => {
        return (
          <CardImage key={i} {...cardData} />
        )
      })}
    </div>
  )
}