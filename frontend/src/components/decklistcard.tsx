import { ApiPostRequest } from "@/requests/requests"
import { OwnedCard } from "@/types/card"
import { Dropdown, DropdownTrigger, DropdownMenu, DropdownItem, dropdown } from "@nextui-org/react"
import { useState } from "react"
import { ManaCost } from "./manacost"

export type DecklistCardProps = {
  deckId?: string
  card: OwnedCard
  board?: string
  count: number
  onHoverFn?: () => void
  refreshDeckFn?: () => void
}

export type DecklistCardListProps = DecklistCardProps & {}

export function DecklistCardList(props: DecklistCardListProps) {
  let [dropdownOpen, setDropdownOpen] = useState<boolean>()

  let moveToBoard = (target_board: string, amount: number) => {
    ApiPostRequest({
      route: "/deck/card/remove",
      body: {
        owned_card_id: props.card.id,
        board: props.board,
        deck_id: props.deckId,
        amount: amount,
      },
      errorHandler: () => { },
      responseHandler: () => {
        target_board != "" ? ApiPostRequest({
          route: "/deck/card",
          body: {
            owned_card_id: props.card.id,
            board: target_board,
            deck_id: props.deckId,
            amount: amount,
          },
          errorHandler: () => { },
          responseHandler: () => {
            if (props.refreshDeckFn) props.refreshDeckFn()
          }
        }) : (props.refreshDeckFn) ? props.refreshDeckFn() : null
      }
    })
  }

  let addOne = () => {
    ApiPostRequest({
      route: "/deck/card",
      body: {
        owned_card_id: props.card.id,
        board: props.board,
        deck_id: props.deckId,
        amount: 1,
      },
      errorHandler: () => { },
      responseHandler: () => { (props.refreshDeckFn) ? props.refreshDeckFn() : null }
    })
  }

  let dropdownItems = [
    {
      key: "add-one",
      label: "Add one copy",
      action: () => addOne(),
      boards: ["b_mainboard", "b_sideboard", "b_maybeboard"],
      isDisabled: () => props.count >= props.card.count
    },
    {
      key: "move-mainboard-one",
      label: "Move one to main deck",
      action: () => moveToBoard("b_mainboard", 1),
      boards: ["b_sideboard", "b_maybeboard"]
    },
    {
      key: "move-mainboard-all",
      label: "Move ALL to main deck",
      action: () => moveToBoard("b_mainboard", props.count),
      boards: ["b_sideboard", "b_maybeboard"]
    },
    {
      key: "move-sideboard-one",
      label: "Move one to sideboard",
      action: () => moveToBoard("b_sideboard", 1),
      boards: ["b_mainboard", "b_maybeboard"]
    },
    {
      key: "move-sideboard-all",
      label: "Move ALL to sideboard",
      action: () => moveToBoard("b_sideboard", props.count),
      boards: ["b_mainboard", "b_maybeboard"]
    },
    {
      key: "move-maybeboard-one",
      label: "Move one to considering",
      action: () => moveToBoard("b_maybeboard", 1),
      boards: ["b_mainboard", "b_sideboard"]
    },
    {
      key: "move-maybeboard-all",
      label: "Move ALL to considering",
      action: () => moveToBoard("b_maybeboard", props.count),
      boards: ["b_mainboard", "b_sideboard"]
    },
    {
      key: "delete-one",
      label: "Remove one",
      action: () => moveToBoard("", 1),
      boards: ["b_mainboard", "b_sideboard", "b_maybeboard"]
    },
    {
      key: "delete-all",
      label: "Remove ALL",
      action: () => moveToBoard("", props.count),
      boards: ["b_mainboard", "b_sideboard", "b_maybeboard"]
    }
  ]

  return (
    <div>
      <Dropdown>
        <DropdownTrigger>
          <div
            className="text-white flex flex-row justify-between items-center cursor-pointer hover:pl-1 duration-75"
            onMouseEnter={props.onHoverFn}
            onClick={() => setDropdownOpen(!dropdownOpen)}>
            <div className="flex flex-row items-start gap-2">
              <p className="font-medium text-gray-400">
                {props.count}
              </p>
              {props.card.card_data.name}
            </div>
            <ManaCost manaCost={props.card.card_data.mana_cost} />
          </div>
        </DropdownTrigger>
        <DropdownMenu
          aria-label="actions"
          items={dropdownItems.filter(
            (item) => item.boards.includes(props.board || "")
          )}
          disabledKeys={
            dropdownItems
              .filter((item) => item.isDisabled ? item.isDisabled() : false)
              .map((item) => item.key)
          }
        >
          {(item) => (
            <DropdownItem
              key={item.key}
              color={item.key === "delete-one" || item.key === "delete-all" ? "danger" : "default"}
              className={item.key === "delete-one" || item.key === "delete-all" ? "text-danger" : "text-white"}
              onClick={item.action}
            >
              {item.label}
            </DropdownItem>
          )}
        </DropdownMenu>
      </Dropdown>
    </div>
  )
}

export type DecklistListProps = {
  category: string
  cards: DecklistCardProps[]
}

export function DecklistList(props: DecklistListProps) {
  if (!props.cards.length) return
  let totalcards = props.cards.map((card) => card.count).reduce((ov, pv) => ov + pv)
  return (
    <div className="flex flex-row">
      <div className="flex flex-col gap-2 px-8 my-4 w-[340px]">
        <div className="flex flex-row justify-between items-center">
          <p className="text-white font-bold">{props.category}</p>
          <p className="text-white">{totalcards}</p>
        </div>
        <div className="h-0.5 bg-white opacity-10" />
        <div className="flex flex-col flex-wrap">
          {props.cards.map((card) => (
            <DecklistCardList key={card.card.id} {...card} />
          ))}
        </div>
      </div >
      <div className="w-[2px] h-full bg-white opacity-10" />
    </div>
  )
}