import './App.css'
import { Navbar01 } from '@/components/ui/shadcn-io/navbar-01';
import Example from './components/ui/example';


function App() {

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar01 className="w-full" />
      <main className="flex-1">
        <Example />
      </main>
    </div>
  )
}

export default App
