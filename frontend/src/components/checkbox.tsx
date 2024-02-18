export type CheckboxProps = React.PropsWithChildren & {
  color?: "primary" | "secondary"
}

export function Checkbox(props: CheckboxProps) {
  return (
    <div className="flex items-start">
      <div className="flex items-center h-5">
        <input id="remember" aria-describedby="remember" type="checkbox" className={`w-4 h-4 border border-gray-300 rounded bg-gray-50 focus:ring-3 focus:ring-${props.color || "primary"}-300`} required={false} />
      </div>
      <div className="ml-3 text-sm">
        <label htmlFor="remember" className="text-gray-400">{props.children}</label>
      </div>
    </div>
  )
}