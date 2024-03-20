"use client"

import { CardData } from "@/types/card"
import Image from "next/image"
import { useEffect, useState } from "react"

export type CardFullProps = {
  card: CardData
  flipped: boolean
  showRarityWhenFlipped: boolean
  onClickFn: () => void
  count?: number
  disabled?: boolean
  extraContent?: React.ReactElement
  inDeck?: number
}

export function CardFull(props: CardFullProps) {
  let borderRarityColor = {
    "common": "border-rarity-common",
    "uncommon": "border-rarity-uncommon",
    "rare": "border-rarity-rare",
    "mythic": "border-rarity-mythic",
  }[props.card.rarity]
  let shadowRarityColor = {
    "common": "shadow-rarity-common",
    "uncommon": "shadow-rarity-uncommon",
    "rare": "shadow-rarity-rare",
    "mythic": "shadow-rarity-mythic",
  }[props.card.rarity]

  let [backImage, setBackImage] = useState<string>("/cardback.webp")

  useEffect(() => {
    if (props.flipped && props.card.back_image_url != "") setBackImage(props.card.back_image_url)
  }, [props.flipped, props.card.back_image_url])

  return (
    <div className="group w-[256px] h-[355px] hover:scale-110 will-change-transform scale-100 duration-75 z-[100] hover:z-[110] [perspective:1000px]">
      {props.disabled && <div className="absolute w-full h-full bg-black opacity-40 z-10" />}
      {props.count && props.count > 1 ? <div className="absolute z-[200] bg-lime-950 border-1 border-white text-white text-lg font-bold px-2 -my-1 -mx-1 rounded">{props.count}</div> : ""}
      <div onClick={() => { if (!props.disabled) props.onClickFn() }} className={
        `absolute rounded-xl w-full h-full duration-500 transition-all [transform-style:preserve-3d] ${!props.flipped && "[transform:rotateY(180deg)]"}`
      }>
        <div className="absolute inset-0 [backface-visibility:hidden]">
          <Image
            className={`duration-75 border-2 ${borderRarityColor} rounded-xl ${!props.flipped && props.showRarityWhenFlipped && "shadow-[0px_0px_20px_1px_rgba(0,0,0,0.3)]"} ${shadowRarityColor}`}
            unoptimized
            priority
            src={props.card.image_url}
            alt={props.card.image_url}
            width={1024}
            height={768}
            quality={100}
          />
        </div>
        <div className="absolute inset-0 [backface-visibility:hidden] [transform:rotateY(180deg)]">
          <Image
            className="duration-75 border-2 border-white rounded-xl"
            src={backImage}
            alt="back"
            width={1024}
            height={1024}
            quality={100}
            layout=""
          />
        </div>
      </div>
      {props.inDeck && props.inDeck > 0 ?
        <div className="absolute p-0.5 px-2 m-1 bottom-0 bg-gray-800 text-white text-sm rounded-full border-1 border-white">
          {props.inDeck} in deck
        </div> : null}
    </div >
  )
}

export type CardDisplaySpoilerProps = {
  cards: CardFullProps[]
}

export function CardDisplaySpoiler(props: CardDisplaySpoilerProps) {
  if (props.cards.length == 0) {
    return <div className="text-white self-center font-thin italic">No cards to show</div>
  }
  return (
    <div className="flex flex-wrap flex-row gap-4 items-center justify-center">
      {props.cards.map((card: CardFullProps, index: number) => {
        return (<CardFull {...card} key={card.card.image_url + index} />)
      })
      }
    </div >
  )
}