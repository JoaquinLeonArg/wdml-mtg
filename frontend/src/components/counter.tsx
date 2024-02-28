import { ChangeEvent, useEffect, useState } from "react"

export type CounterInputProps = {
  label: string
  id: string
  min: number
  max: number
  initial: number
  onChange: (value: number) => void
}

export function CounterInput(props: CounterInputProps) {
  let [value, setValue] = useState<number>(props.initial)
  useEffect(() => props.onChange(value), [value])

  return (
    <>
      <label htmlFor="quantity-input" className="block mb-2 text-sm font-medium text-white">{props.label}</label>
      <div className="relative flex items-center max-w-[6rem]">
        <button onClick={() => { if (value > props.min) setValue(value - 1) }} type="button" id="decrement-button" data-input-counter-decrement="quantity-input" className="bg-gray-700 hover:bg-gray-600 border-gray-600 border rounded-s-lg p-2 h-8 focus:ring-gray-700 focus:ring-2 focus:outline-none">
          <svg className="w-3 h-3 text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 18 2">
            <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M1 1h16" />
          </svg>
        </button>
        <input type="text" id={props.id} data-input-counter data-input-counter-min={props.min} data-input-counter-max={props.max} className="border-x-0  h-8 text-center text-sm block w-full py-2.5 bg-gray-700 border-gray-600 placeholder-gray-400 text-white focus:ring-blue-500 focus:border-blue-500" value={value} disabled required />
        <button onClick={() => { if (value < props.max) setValue(value + 1) }} type="button" id="increment-button" data-input-counter-increment="quantity-input" className="bg-gray-700 hover:bg-gray-600 border-gray-600 border rounded-e-lg p-2 h-8 focus:ring-gray-700 focus:ring-2 focus:outline-none">
          <svg className="w-3 h-3 text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 18 18">
            <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 1v16M1 9h16" />
          </svg>
        </button>
      </div>
    </>
  )
}