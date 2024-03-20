import React from "react"

export type ManaCostProps = {
  manaCost: string
}

export function ManaCost(props: ManaCostProps) {
  let re = /{([0-9A-Z\/]*)}/g
  let childs: React.ReactElement[] = []
  for (const symbols of props.manaCost.matchAll(re)) {
    childs.push(<i className={`ms ms-${symbols[1].split("/").join("").toLowerCase()} ms-cost ms-shadow`}></i>)
  }

  return (
    <div className="flex flex-row gap-0.5">
      {childs}
    </div>
  )
}