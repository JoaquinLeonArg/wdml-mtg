import Combobox from "react-widgets/Combobox";

export type ComboBoxProps = {
  default: string
  data: string[]
}

export function ComboBox(props: ComboBoxProps) {
  return (
    <Combobox
      defaultValue={props.default}
      data={props.data}
    // className=""
    // containerClassName=""
    />
  )
}