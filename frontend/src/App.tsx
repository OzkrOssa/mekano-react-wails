import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import Payment from "@/components/Payment";
import Billing from "@/components/Billing";

function App() {
  return (
    <div className="flex justify-center items-center m-5">
    <Tabs defaultValue="payment" className="w-[400px]">
      <TabsList className="grid w-full grid-cols-2">
        <TabsTrigger value="payment">Pagos</TabsTrigger>
        <TabsTrigger value="billing">Facturacion</TabsTrigger>
      </TabsList>
      <TabsContent value="payment">
        <Payment />
      </TabsContent>
      <TabsContent value="billing">
        <Billing />
      </TabsContent>
    </Tabs>
    </div>
  );
}

export default App;
