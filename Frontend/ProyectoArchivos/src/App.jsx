function App() {
  return (
    <div className="max-w-md mx-auto">
      <form className="bg-slate-800 p-10 mb-4">
        <input
          className="bg-slate-300 p-3 w-full mb-2"
          placeholder=""
          autoFocus
        />
        <textarea
          className="bg-slate-300 p-3 w-full mb-2"
          placeholder=""
        ></textarea>
        <button className="bg-blue-600 hover:bg-blue-400 text-yellow-300">Ejecutar</button>
      </form>
    </div>
  );
}
export default App;
