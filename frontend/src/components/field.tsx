import { ChangeEvent } from "react"

export type TextFieldWithLabelProps = {
  id: string
  label: string
  type?: string
  color?: "primary" | "secondary"
  htmlFor?: string
  placeholder?: string
  size?: number
  required?: boolean
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void
}

export function TextFieldWithLabel(props: TextFieldWithLabelProps) {
  return (
    <div>
      <label htmlFor={props.htmlFor || ""} className="block mb-2 text-sm font-medium text-gray-100">{props.label}</label>
      <input
        required={props.required || false}
        onChange={props.onChange}
        size={props.size || 30}
        type={props.type || ""}
        name={props.id}
        id={props.id}
        className={`bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-${props.color || "primary"}-600 focus:border-${props.color || "primary"}-600 block w-full p-2.5`}
        placeholder={props.placeholder || ""} />
    </div>
  )
}