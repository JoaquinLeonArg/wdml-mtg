
export type PrimaryButtonProps = React.PropsWithChildren & {
  icon?: "arrow"
  type?: "primary" | "secondary" | "cancel"
  fullWidth?: boolean
}

export function Button(props: PrimaryButtonProps) {
  let bgClasses = {
    "primary": "bg-primary-700",
    "secondary": "bg-secondary-700",
    "cancel": "bg-rose-700"
  }
  return (
    <a href="#" className={`${props.fullWidth && "w-full"} inline-flex items-center justify-center text-white ${bgClasses[props.type || "primary"]} hover:bg-primary-800 focus:ring-4 focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:focus:ring-primary-900`} >
      {props.children}
      {props.icon == "arrow" && <svg className="ml-2 -mr-1 w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>}
    </a >
  )
}