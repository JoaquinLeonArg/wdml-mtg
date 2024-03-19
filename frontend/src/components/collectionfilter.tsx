import { Input, Dropdown, DropdownTrigger, Button, DropdownMenu, DropdownItem, ButtonGroup, colors } from "@nextui-org/react"
import { useState } from "react"

export type CollectionFilterProps = {
  count?: number
  setCardName: (cardName: string) => void
  setTags: (tags: string) => void
  setRarity: (rarity: keyof typeof rarities) => void
  rarity: keyof typeof rarities
  setColors: (cols: mtgColors) => void
  colors: mtgColors
  setTypes: (types: string) => void
  setOracle: (oracle: string) => void
  setSetCode: (setCode: string) => void
  setMv: (mv: string) => void
}

export const rarities = {
  "": "Any rarity",
  "common": "Common",
  "uncommon": "Uncommon",
  "rare": "Rare",
  "mythic": "Mythic"
}

export type mtgColors = { W: boolean, U: boolean, B: boolean, R: boolean, G: boolean, C: boolean }

export function CollectionFilter(props: CollectionFilterProps) {
  return (
    <div>
      <div className="flex flex-col w-full">
        <div className="flex flex-row gap-2 mb-2 items-center">
          <Input
            onChange={(e) => props.setCardName(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Card name"
            className="text-white max-w-96"
          />
          <Input
            onChange={(e) => props.setTags(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Tags"
            className="text-white max-w-96"
          />
          <Dropdown>
            <DropdownTrigger>
              <Button
                variant="bordered"
                className="capitalize"
              >
                {rarities[props.rarity]}
              </Button>
            </DropdownTrigger>
            <DropdownMenu
              className="text-white"
            >
              {
                Object.keys(rarities).map((key) =>
                  <DropdownItem onPress={() => props.setRarity(key as keyof typeof rarities)} key={key}>{rarities[key as keyof typeof rarities]}</DropdownItem>
                )
              }
            </DropdownMenu>
          </Dropdown>
          <ButtonGroup variant="ghost" fullWidth={false}>
            {
              Object.keys(props.colors).map((key: string) =>
                <Button isIconOnly
                  className={`text-xl ${colors[key as keyof typeof colors] && "bg-gray-500"}`}
                  variant="ghost"
                  onClick={() => props.setColors({ ...props.colors, [key]: !props.colors[key as keyof mtgColors] })}
                  key={key}
                >
                  <i className={`ms ms-${key.toLowerCase()} ms-cost`}></i>
                </Button>
              )
            }
          </ButtonGroup>
        </div>
        <div className="flex flex-row gap-2 mb-2 items-center">
          <Input
            onChange={(e) => props.setTypes(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Types"
            className="text-white max-w-96"
          />
          <Input
            onChange={(e) => props.setOracle(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Oracle"
            className="text-white max-w-96"
          />
          <Input
            onChange={(e) => props.setSetCode(e.target.value)}
            variant="bordered"
            type="string"
            placeholder="Set code"
            className="text-white max-w-24"
          />
          <Input
            onChange={(e) => props.setMv(e.target.value)}
            variant="bordered"
            type="number"
            placeholder="Any"
            className="text-white max-w-28"
            endContent={
              <div className="pointer-events-none flex items-center">
                <span className="text-gray-300 text-small">MV</span>
              </div>
            }
          />
        </div>
        {props.count && (
          <div className="text-white italic font-thin">
            Found: {props.count} cards
          </div>
        )}
      </div>
    </div>
  )
}