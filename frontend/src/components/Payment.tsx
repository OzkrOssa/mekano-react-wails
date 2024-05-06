import React from "react";
import ReactLoading from "react-loading";
import {
  Card,
  CardContent,
  CardDescription,
  CardTitle,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Upload } from "lucide-react";
import { OpenFile, MekanoPayment } from "@/lib/wailsjs/go/main/App";
import { mekano } from "@/lib/wailsjs/go/models";

export default function Payment() {
  const [payFilePath, setPayFilePath] = React.useState("");
  const [processing, setProcessing] = React.useState(false);
  const [result, setResult] = React.useState<mekano.PaymentStatistics | null>(null);

  const OpenDialog = () => {
    OpenFile().then((path) => {
      setPayFilePath(path);
    });
  };

  const ProcessPayment = () => {
    setProcessing(true);
    MekanoPayment(payFilePath).then((result) => {
      setResult(result);
      setProcessing(false);
    });
  }
  return (
    <Card>
      <CardHeader>
        <CardTitle>Pagos</CardTitle>
        <CardDescription>Carga el archivo de pagos aqui.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-2">
        <div className="relative">
          <Upload
            size={35}
            className="absolute flex items-center justify-center pl-3 bg-transparent text-black"
            onClick={OpenDialog}
          />

          <Input
            placeholder={payFilePath}
            className="pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </CardContent>
      <CardFooter>
        <Button
          className="w-full"
          onClick={ProcessPayment}
          disabled={processing}
        >
          {processing ? (
            <ReactLoading type={"spin"} height={20} width={20} />
          ) : (
            "Procesar"
          )}
        </Button>
      </CardFooter>
    </Card>
  );
}
