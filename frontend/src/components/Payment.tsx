import React from "react";
import ReactLoading from "react-loading";
import * as AlertDialog from "@/components/ui/alert-dialog";
import * as Table from "@/components/ui/table";
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

  const handleClear = () => {
    setPayFilePath("");
    setResult(null);
  }
  return (
    <React.Fragment>
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
    {result !== null ? (
        <AlertDialog.AlertDialog defaultOpen={true}>
          <AlertDialog.AlertDialogContent>
            <AlertDialog.AlertDialogHeader>
              <AlertDialog.AlertDialogTitle>
                Resultados
              </AlertDialog.AlertDialogTitle>
              <AlertDialog.AlertDialogDescription>
                <Table.Table>
                  <Table.TableCaption>{result.archivo}</Table.TableCaption>
                  <Table.TableCaption className="font-bold text-sm">
                    Rango RC: {result.rango_rc}
                  </Table.TableCaption>
                  <Table.TableHeader>
                    <Table.TableRow>
                      <Table.TableHead className="w-[100px] font-bold">
                        Caja
                      </Table.TableHead>
                      <Table.TableHead className="text-right font-bold">
                        Valor
                      </Table.TableHead>
                    </Table.TableRow>
                  </Table.TableHeader>
                  <Table.TableBody>
                    <Table.TableRow>
                      <Table.TableCell className="font-medium">
                        Bancolombia
                      </Table.TableCell>
                      <Table.TableCell className="text-right">
                        ${new Intl.NumberFormat().format(result.bancolombia)}
                      </Table.TableCell>
                    </Table.TableRow>
                    <Table.TableRow>
                      <Table.TableCell className="font-medium">
                        Davivienda
                      </Table.TableCell>
                      <Table.TableCell className="text-right">
                        ${new Intl.NumberFormat().format(result.davivienda)}
                      </Table.TableCell>
                    </Table.TableRow>
                    <Table.TableRow>
                      <Table.TableCell className="font-medium">
                        Pay U
                      </Table.TableCell>
                      <Table.TableCell className="text-right">
                        ${new Intl.NumberFormat().format(result.payu)}
                      </Table.TableCell>
                    </Table.TableRow>
                    <Table.TableRow>
                      <Table.TableCell className="font-medium">
                        Susuerte
                      </Table.TableCell>
                      <Table.TableCell className="text-right">
                        ${new Intl.NumberFormat().format(result.susuerte)}
                      </Table.TableCell>
                    </Table.TableRow>
                    <Table.TableRow>
                      <Table.TableCell className="font-medium">
                        Efectivo
                      </Table.TableCell>
                      <Table.TableCell className="text-right">
                        ${new Intl.NumberFormat().format(result.efectivo)}
                      </Table.TableCell>
                    </Table.TableRow>
                    <Table.TableRow>
                      <Table.TableCell className="font-medium">
                        Total
                      </Table.TableCell>
                      <Table.TableCell className="text-right">
                        ${new Intl.NumberFormat().format(result.total)}
                      </Table.TableCell>
                    </Table.TableRow>
                  </Table.TableBody>
                </Table.Table>
              </AlertDialog.AlertDialogDescription>
            </AlertDialog.AlertDialogHeader>
            <AlertDialog.AlertDialogFooter>
              <AlertDialog.AlertDialogAction onClick={handleClear}>
                Continue
              </AlertDialog.AlertDialogAction>
            </AlertDialog.AlertDialogFooter>
          </AlertDialog.AlertDialogContent>
        </AlertDialog.AlertDialog>
      ) : null}
    </React.Fragment>
  );
}
